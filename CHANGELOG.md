YOU CANNOT DOWNLOAD THIS UPDATE THROUGH MAGISK'S IN-APP MODULE UPDATER.

Install this release here: [Release v2.0.2](https://github.com/JoshuaDoes/pixel-tensor-audio-decompressor/releases/tag/202)

### v2.0.2 (it's 2:02 AM, and Magisk is silly once again)
- Break the ability to update the module through Magisk
- For some reason, ZIPs downloaded through the Magisk in-app module updater will extract without executable permissions
- We cannot have this, or the menu will cease to function
- Until this issue can be resolved on Magisk's end, you must download all future updates from the GitHub Releases page
- You will still be notified of future updates, but they will fail to download from now on until this issue is resolved

### v2.0.0 (menu debut)
- Integrate my custom text menu engine! Use volume keys to navigate, and tap the screen to select a menu item
- Support live testing of different volume presets before committing to a choice, but it will abort the install! You must reinstall to demo another preset
- Require navigation to the install patch option, to prevent accidentally installing a choice you may not like
- Recommended volume is still 865/14, but may change with future updates depending on community feedback
- Allow cancelling the install before choosing other options

The current volume presets you may choose from:

- 865/14: The recommended volume by my testing
- 865/13: The original volume from v1.0.0
- 865/12: Recommended by a Pixel 7 Pro (cheetah) user
- 817/17: The Google stock volume
- 913/15: Potentially unsafe! Triggers volume normalization when the volume spikes
- 913/20: Potentially even more unsafe than 913/15! The absolute maximum that can be set

More presets may be added in future updates, also depending on community feedback

### v1.2.0 (tangorfelix)
- Add experimental support for the Pixel Fold (felix) and the Pixel Tablet (tangorpro)
- Allow installation on Android 14

### v1.1.2 (Magisk is silly)
- Hotfix to remove the `skip_mount` file after the install process... lol

### v1.1.1 (oops)
- Hotfix update to fix the install script, sorry!

### v1.1.0 (we can do better)
- Raised the amplifier PCM gain from 13 to 14. Slightly increases the perceived volume, making the bass a little fuller
- Modified the sed function of the install script to work regardless of current values

### v1.0.0 (initial release)
- Raised digital PCM volume from stock value of 817 to 865
- Lowered amplifier PCM gain from stock value of 17 to 13
- Forced boost peak current limit to 3.50A. Default on most ROMs including Google stock, but no longer default for hentaiOS and its forks, potentially others
