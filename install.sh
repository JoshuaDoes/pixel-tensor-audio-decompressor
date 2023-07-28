SKIPMOUNT=false

print_modname() {
  ui_print ""
  ui_print "***********************************"
  ui_print "* Pixel Tensor Audio Decompressor *"
  ui_print "             * 1.2.0 *             "
  ui_print "** Made and tested by JoshuaDoes **"
  ui_print "***********************************"
  ui_print "For the following devices:"
  ui_print "- Pixel 6 Pro (raven)"
  ui_print "- Pixel 6     (oriole)"
  ui_print "- Pixel 6a    (bluejay)"
  ui_print "- Pixel 7 Pro (cheetah)"
  ui_print "- Pixel 7     (panther)"
  ui_print "- Pixel 7a    (lynx)"
  ui_print "In testing for these devices:"
  ui_print " - Pixel Fold   (felix)"
  ui_print " - Pixel Tablet (tangorpro)"
  ui_print "***********************************"
  ui_print "The following changes will be made!"
  ui_print ""
  ui_print "# Speakers"
  ui_print "Digital PCM volume:       817 -> 865"
  ui_print "Amplifier PCM gain:       17 -> 14"
  ui_print "Boost peak current limit: 3.50A"
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

  if [ $DEVICE = "felix" ] || [ $DEVICE = "tangorpro" ]; then
    ui_print "* WARNING: "$DEVICE" is experimental!"
  fi

  if [ $RELEASE != "13" ] && [ $RELEASE != "14" ]; then
    abort "* Android "$RELEASE" is not supported!"
  fi

  ui_print "- Copying mixer paths to a temporary location"
  cp "/vendor/etc/mixer_paths.xml" "$TMPDIR/mixer_paths.xml"
  ui_print "- Patching mixer paths using sed"
  sed -r 's/AMP PCM Gain" value="(.*)"/AMP PCM Gain" value="14"/g; s/Digital PCM Volume" value="(.*)"/Digital PCM Volume" value="865"/g; s/Boost Peak Current Limit" value="(.*)A"/Boost Peak Current Limit" value="3.50A"/g' "$TMPDIR/mixer_paths.xml" > "$TMPDIR/mixer_paths2.xml"
  ui_print "- Installing patched mixer paths systemfully into module path"
  mkdir -p "$MODPATH/system/vendor/etc"
  mv "$TMPDIR/mixer_paths2.xml" "$MODPATH/system/vendor/etc/mixer_paths.xml"
  ui_print "- Cleaning up temporary files"
  rm "$TMPDIR/mixer_paths.xml"
  rm "$TMPDIR/mixer_paths2.xml"

  ui_print ""
  ui_print "* Reboot and listen to some audio!"
  ui_print ""
}

set_permissions() {
  set_perm_recursive $MODPATH 0 0 0755 0644
}
