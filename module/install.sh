SKIPMOUNT=false

print_modname() {
  ui_print ""
  ui_print "***********************************"
  ui_print "* Pixel Tensor Audio Decompressor *"
  ui_print "             * 2.1.0 *             "
  ui_print "** Made and tested by JoshuaDoes **"
  ui_print "***********************************"
  ui_print "For the following devices:"
  ui_print "- Pixel 6 Pro  (raven)"
  ui_print "- Pixel 6      (oriole)"
  ui_print "- Pixel 6a     (bluejay)"
  ui_print "- Pixel 7 Pro  (cheetah)"
  ui_print "- Pixel 7      (panther)"
  ui_print "- Pixel 7a     (lynx)"
  ui_print "- Pixel Fold   (felix)"
  ui_print "- Pixel Tablet (tangorpro)"
  ui_print "***********************************"
  ui_print ""
}

on_install() {
  # raven = Pixel 6 Pro
  # oriole = Pixel 6
  # bluejay = Pixel 6a
  # cheetah = Pixel 7 Pro
  # panther = Pixel 7
  # lynx = Pixel 7a
  # felix = Pixel Fold
  # tangorpro = Pixel Tablet
  if [ $DEVICE != "raven" ] && [ $DEVICE != "oriole" ] && [ $DEVICE != "bluejay" ] && [ $DEVICE != "cheetah" ] && [ $DEVICE != "panther" ] && [ $DEVICE != "lynx" ] && [ $DEVICE != "felix" ] && [ $DEVICE != "tangorpro" ]; then
    abort "* "$DEVICE" is not supported!"
  fi

  if [ $RELEASE != "13" ] && [ $RELEASE != "14" ]; then
    abort "* Android "$RELEASE" needs testing!"
  fi

  ui_print "If you don't see the menu,"
  ui_print "close Magisk and try again."
  ui_print ""
  ui_print "The official Magisk app will"
  ui_print "provide the best experience."
  ui_print ""
  ui_print "AMM also handles the menu well."
  ui_print "KernelSU does not like it. :c"
  sleep 1

  export TERM=xterm-256color
  mkdir "$MODPATH"
  unzip -o "$ZIPFILE" "*" -d "$MODPATH" >&2
  exec "$MODPATH/bin/menu" --workingDir "$MODPATH" --keyCalibration "/data/local/tmp/menuKeycodes.json" 2>&1
}

set_permissions() {
  set_perm_recursive $MODPATH 0 0 0755 0644
}
