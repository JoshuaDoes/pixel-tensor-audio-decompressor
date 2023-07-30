package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// MenuItem holds an item for a menu, such as a button, a checkbox, or an input box
type MenuItem struct {
	Text   string `json:"text"`
	Desc   string `json:"desc"`
	Type   string `json:"type"`   //menu, exec, explorer[:pwd], note, var name
	Action string `json:"action"` //var: string[:limit]|number[:min[:max]]|file[:extension1[,extension2,...]]|bool|opts:opt1,opt2,[opt3,...]
}

// MenuItemList holds a list of items to interact with
type MenuItemList struct {
	Title      string      `json:"title"`
	Subtitle   string      `json:"subtitle"`
	Items      []*MenuItem `json:"items"`      //items to display on the page
	NoGoBack   bool        `json:"noGoBack"`   //hides the go back button
	NoSelector bool        `json:"noSelector"` //hides the item cursor
	DefaultCur int         `json:"defaultCur"` //the cursor to set by default
	Exec       string      `json:"exec"`       //a line interpreted as an exec action
}

func (m *MenuItemList) AddItem(name, desc, itemType, action string) {
	m.Items = append(m.Items, &MenuItem{Text: name, Desc: desc, Type: itemType, Action: action})
}

// MenuEngine holds a list of menus and acts as the menu interface
type MenuEngine struct {
	//Menu navigation
	Menus       map[string]*MenuItemList
	HomeMenu    string
	LoadedMenu  string
	MenuHistory []string
	ItemHistory []int
	Environment map[string]string //global variables set by menus
	ItemCursor  int
	Locked      bool
	Return      string                          //return value set by some menu types
	Hooks       map[string]func(me *MenuEngine) //run a hook after changing to a menu

	//Rendering control
	Render func(string)
	LinesH int
	LinesV int
}

// NewMenuEngine returns a menu engine ready to be used
func NewMenuEngine(renderer func(string), width, height int) *MenuEngine {
	return &MenuEngine{
		Menus:       make(map[string]*MenuItemList),
		MenuHistory: make([]string, 0),
		ItemHistory: make([]int, 0),
		Environment: make(map[string]string),
		Render:      renderer,
		LinesH:      width,
		LinesV:      height,
		Hooks:       make(map[string]func(me *MenuEngine)),
	}
}

func (me *MenuEngine) Hook(id string, hook func(me *MenuEngine)) {
	me.Hooks[id] = hook
}

func (me *MenuEngine) LoadMenu(id string, itemList *MenuItemList) {
	me.Menus[id] = itemList
}

func (me *MenuEngine) Lock() {
	me.Locked = true
}
func (me *MenuEngine) Unlock() {
	me.Locked = false
}

func (me *MenuEngine) init() {
	if me.Menus == nil {
		me.Menus = make(map[string]*MenuItemList)
	}
	if me.MenuHistory == nil {
		me.MenuHistory = make([]string, 0)
	}
	if me.ItemHistory == nil {
		me.ItemHistory = make([]int, 0)
	}
	if me.Environment == nil {
		me.Environment = make(map[string]string)
	}
	if me.Hooks == nil {
		me.Hooks = make(map[string]func(me *MenuEngine))
	}
}

func (me *MenuEngine) isBackVisible() bool {
	if me.Menus[me.LoadedMenu].NoGoBack {
		return false
	}
	return len(me.MenuHistory) > 0
}

// PrevItem navigates to the previous menu item, or to the last if none previous
func (me *MenuEngine) PrevItem() {
	if me.Locked {
		return
	}
	me.init()
	defer me.render()

	if me.isBackVisible() && me.ItemCursor == -1 {
		me.ItemCursor = len(me.Menus[me.LoadedMenu].Items) - 1
	} else if !me.isBackVisible() && me.ItemCursor == 0 {
		me.ItemCursor = len(me.Menus[me.LoadedMenu].Items) - 1
	} else {
		me.ItemCursor--
	}

	if me.ItemCursor >= 0 && me.Menus[me.LoadedMenu].Items[me.ItemCursor].Type == "divider" {
		me.PrevItem()
	}
}

// NextItem navigates to the next menu item, or to the first if none next
func (me *MenuEngine) NextItem() {
	if me.Locked {
		return
	}
	me.init()
	defer me.render()

	if (me.ItemCursor + 1) >= len(me.Menus[me.LoadedMenu].Items) {
		if me.isBackVisible() {
			me.ItemCursor = -1
		} else {
			me.ItemCursor = 0
		}
	} else {
		me.ItemCursor++
	}

	if me.ItemCursor >= 0 && me.Menus[me.LoadedMenu].Items[me.ItemCursor].Type == "divider" {
		me.NextItem()
	}
}

// Action activates the selected item's action, such as navigating to a menu or executing a program
func (me *MenuEngine) Action() {
	if me.Locked {
		return
	}
	me.init()

	if me.ItemCursor == -1 {
		me.PrevMenu()
		return
	}

	if len(me.Menus[me.LoadedMenu].Items) == 0 {
		return
	}

	selectedItem := me.Menus[me.LoadedMenu].Items[me.ItemCursor]
	selectedAction := me.Vars(selectedItem.Action)
	itemArgs := strings.Split(me.Vars(selectedItem.Type), " ")
	actionArgs := strings.Split(selectedAction, " ")
	switch itemArgs[0] {
	case "internal":
		switch actionArgs[0] {
		case "abort":
			if len(actionArgs) > 1 {
				me.ChangeMenu(actionArgs[1])
			}
			os.Exit(1)
		case "exit":
			if len(actionArgs) > 1 {
				me.ChangeMenu(actionArgs[1])
			}
			os.Exit(0)
		default:
			me.ErrorText("Unknown internal action", selectedAction)
		}
	case "menu":
		me.ChangeMenu(actionArgs[0])
	case "exec":
		me.RunRealtime(selectedAction)
	case "explorer":
		workingDir := "/"
		if len(itemArgs) > 1 {
			workingDir = strings.Join(itemArgs[1:], " ")
		}
		me.Explorer(workingDir, selectedAction)
	case "return":
		if me.Return != "" {
			me.Environment[me.Return] = selectedAction
			me.Return = ""
		}
		me.PrevMenu()

		//Back all the way out of an explorer context
		for {
			if string(me.Menus[me.LoadedMenu].Title[:8]) == "Explorer" {
				me.PrevMenu()
				continue
			}
			break
		}
	case "setvar":
		me.Return = itemArgs[1] //set var for what to return to
		if len(itemArgs) > 2 {
			me.Environment[itemArgs[1]] = itemArgs[2] //The value was supplied by the menu
		}
		for i := 3; i < len(itemArgs); i++ {
			if i+1 < len(itemArgs) {
				me.Environment[itemArgs[i]] = itemArgs[i+1]
			}
			i++
		}

		switch actionArgs[0] {
		case "explorer":
			workingDir := "/"
			if len(actionArgs) > 1 {
				workingDir = strings.Join(actionArgs[1:], " ")
			}
			me.Explorer(workingDir, "")
		case "menu":
			me.ChangeMenu(actionArgs[1])
		default:
			me.ErrorText("Unknown action for var " + me.Return, selectedAction)
		}
	case "note":
		if selectedAction != "" {
			me.DisplayText(selectedAction)
		} else {
			me.Redraw() //hide the newline
		}
	default:
		me.ErrorText("Unknown action: " + selectedItem.Type, selectedAction)
	}
}

// Explorer abuses the powers of AddMenu, ChangeMenu, and PrevMenu to create a file browser with support for passing a selected file to an executable
func (me *MenuEngine) Explorer(workingDir, bin string) {
	if workingDir[len(workingDir)-1] != '/' {
		workingDir += "/"
	}

	displayBin := workingDir
	if bin != "" {
		displayBin = strings.Replace(bin, "$?", workingDir, -1)
	}
	explorer := &MenuItemList{
		Title: "Explorer - " + displayBin,
		Items: make([]*MenuItem, 0),
	}

	files, err := ioutil.ReadDir(workingDir)
	if err != nil {
		explorer.AddItem("Failed to list the files in "+workingDir, fmt.Sprintf("%v", err), "note", "")
	} else {
		for _, file := range files {
			if file.IsDir() {
				explorer.AddItem(file.Name()+"/", workingDir+file.Name()+"/", "explorer "+workingDir+file.Name(), bin)
			} else {
				if bin != "" {
					explorer.AddItem(file.Name(), "", "exec", strings.Replace(bin, "$?", fmt.Sprintf("%s%s", workingDir, file.Name()), -1))
				} else {
					explorer.AddItem(file.Name(), "", "return", workingDir+file.Name())
				}
			}
		}
	}

	me.AddMenu(workingDir, explorer)
	me.ChangeMenu(workingDir)
}

// RunRealtime runs the given command, but doesn't halt the menu engine
func (me *MenuEngine) RunRealtime(command string) {
	me.Lock()
	defer me.Unlock()
	err := RunRealtime(me.Vars(command))
	if err != nil {
		me.ErrorText(err.Error(), "")
		return
	}
}

// Run runs the given command, but halts the menu engine until completion
func (me *MenuEngine) Run(command string) {
	me.Lock()
	defer me.Unlock()
	out, err := Run(me.Vars(command))
	outString := string(out)
	if err != nil {
		me.ErrorText(err.Error(), outString)
		return
	}
	me.Menus[me.LoadedMenu].AddItem(outString, "Task complete", "note", "")
	me.render()
}

// AddMenu adds a menu to the menu list
func (me *MenuEngine) AddMenu(menuID string, menu *MenuItemList) {
	me.init()
	me.Menus[menuID] = menu
}

// RemoveMenu removes a menu from the menu list
func (me *MenuEngine) RemoveMenu(menuID string) {
	me.init()
	me.Menus[menuID] = nil
}

// ChangeMenu changes to another available menu
func (me *MenuEngine) ChangeMenu(menuID string) {
	me.init()

	lm, ok := me.Menus[menuID]
	if !ok {
		me.ErrorText("Unknown menu", menuID)
		return
	}

	if me.LoadedMenu != "" { //&& me.LoadedMenu != "INTERNAL_ERROR_TEXT" {
		me.MenuHistory = append(me.MenuHistory, me.LoadedMenu)
		me.ItemHistory = append(me.ItemHistory, me.ItemCursor)
	}

	me.LoadedMenu = menuID
	me.ItemCursor = lm.DefaultCur

	if lm.Exec != "" {
		me.Run(lm.Exec)
	}

	me.render()

	_, ok = me.Hooks[menuID]
	if ok {
		me.Hooks[menuID](me)
	}
}

// Home returns to the home menu
func (me *MenuEngine) Home() {
	me.ChangeMenu(me.HomeMenu)
}

// PrevMenu returns to the last menu in history
func (me *MenuEngine) PrevMenu() {
	me.init()
	defer me.render()

	if len(me.MenuHistory) == 0 {
		return //We can't go back to nothing, or can we?
	}

	menuID := me.MenuHistory[len(me.MenuHistory)-1]         //Get the previous menu
	me.MenuHistory = me.MenuHistory[:len(me.MenuHistory)-1] //Remove this menu from history regardless of it being valid
	itemCursor := me.ItemHistory[len(me.ItemHistory)-1]     //Get the previous item cursor
	me.ItemHistory = me.ItemHistory[:len(me.ItemHistory)-1] //Remove this item cursor from history regardless of it being valid

	_, ok := me.Menus[menuID]
	if !ok {
		//Allow returning to a working menu
		me.MenuHistory = append(me.MenuHistory, me.LoadedMenu)
		me.ItemHistory = append(me.ItemHistory, me.ItemCursor)

		me.ErrorText("Unknown menu", menuID)
		return
	}

	//Reset the item cursor if it's out of bounds
	if itemCursor >= len(me.Menus[menuID].Items) {
		itemCursor = 0
	}

	me.LoadedMenu = menuID
	me.ItemCursor = itemCursor

	_, ok = me.Hooks[menuID]
	if ok {
		me.Hooks[menuID](me)
	}
}

// ErrorText generates an error message menu with menuID "INTERNAL_ERROR_TEXT" and navigates to it
// It is used internally as well as being made available, so refrain from using menuIDs starting with "INTERNAL"
func (me *MenuEngine) ErrorText(err, extra string) {
	menuError := &MenuItemList{
		NoGoBack: true,
		Title:    err,
		Subtitle: extra,
		Items: []*MenuItem{
			{
				Text:   "Return home to start over",
				Desc:   "Because you encountered this error, we can't risk going back",
				Type:   "menu",
				Action: me.HomeMenu,
			},
		},
	}
	me.Menus["INTERNAL_ERROR_TEXT"] = menuError
	me.ChangeMenu("INTERNAL_ERROR_TEXT")
	me.ResetHistory()
}

// DisplayText generates a text message menu with menuID "INTERNAL_DISPLAY_TEXT" and navigates to it
// It is used internally as well as being made available, so refrain from using menuIDs starting with "INTERNAL"
func (me *MenuEngine) DisplayText(txt string) {
	menuTxt := &MenuItemList{
		Title: txt,
	}
	me.Menus["INTERNAL_DISPLAY_TEXT"] = menuTxt
	me.ChangeMenu("INTERNAL_DISPLAY_TEXT")
}

// ResetHistory clears the linked item and menu histories, and resets the cursor to item 0
func (me *MenuEngine) ResetHistory() {
	me.ItemCursor = 0
	me.ItemHistory = make([]int, 0)
	me.MenuHistory = make([]string, 0)
}

// GetRender returns a rendered menu text to be displayed immediately, as the menu state can change freely before and after
func (me *MenuEngine) GetRender() string {
	menu := ""

	lm := me.Menus[me.LoadedMenu]
	if lm.Title != "" {
		menu += lm.Title + "\n\n"
	}
	if lm.Subtitle != "" {
		menu += lm.Subtitle + "\n\n"
	}
	if len(lm.Items) > 0 {
		for i := 0; i < len(lm.Items); i++ {
			switch lm.Items[i].Type {
			case "divider":
				if length, err := strconv.Atoi(lm.Items[i].Action); err != nil {
					for j := 0; j < length; j++ {
						menu += "\n"
					}
				} else {
					menu += "\n"
				}
			default:
				if !lm.NoSelector && me.ItemCursor == i {
					menu += "\t--> "
				} else {
					menu += "\t   "
				}
				menu += lm.Items[i].Text
				if lm.Items[i].Type == "menu" {
					menu += " ..."
				}
				menu += "\n"
			}
		}
	}
	if me.isBackVisible() {
		if me.ItemCursor == -1 {
			menu += "\n\t--> "
		} else {
			menu += "\n\t   "
		}
		menu += "Go back\n\n"
	}
	menu += "\n"
	if !lm.NoSelector {
		if me.ItemCursor < 0 {
			menu += " - Return to the previous menu"
		} else if len(lm.Items) > 0 && lm.Items[me.ItemCursor].Desc != "" {
			menu += " - " + lm.Items[me.ItemCursor].Desc
		}
	}

	return me.Vars(menu)
}

// Vars returns a string formatted with all vars replaced, in order from longest var name to shortest to avoid partial var name replacements
func (me *MenuEngine) Vars(in string) string {
	longest := 0
	vars := make(map[int][]string)
	for varName := range me.Environment {
		varLen := len(varName)
		if varLen > longest {
			longest = varLen
		}
		if _, ok := vars[varLen]; !ok {
			vars[varLen] = make([]string, 0)
		}
		vars[varLen] = append(vars[varLen], varName)
	}
	for i := longest; i > 0; i-- {
		if varNames, ok := vars[i]; ok {
			for _, varName := range varNames {
				in = strings.Replace(in, "$"+varName, me.Environment[varName], -1)
			}
		}
	}
	return in
}

func (me *MenuEngine) Redraw() {
	me.render()
}
func (me *MenuEngine) render() {
	if me.Render != nil {
		me.Render(me.GetRender())
	}
}
