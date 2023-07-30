## WARNINGS!

### It goes without saying, your decision to use this module comes with risks. For example, playing at high volumes can lead to increased heat production, and may be unfavorable in environmental conditions most affected by our continuously heating climate. You lost your warranty the moment you chose to root your phone, and I am not to be held responsible if this module causes speaker failure, battery drain, or worse. I have been successfully using my chosen values of 865/14 for weeks now without any issues, and adapted to the speaker protection using 913/15 for months prior to seeking out balanced values. My speakers are still fine, yours should be fine too.

### If you wish to use this module, please provide feedback! It will be most helpful if you can provide as many details as possible, such as the differences you notice between the stock values and my new values (just toggle the module!), the environmental conditions you are in, if you are noticing increased heating compared to before, etc. The goal is to make sure these are truly the best values, and once confirmed, to submit these new values to as many custom ROMs as possible for their new default values.

### In the meantime, dear custom ROM developers: Please, refrain from adopting these new values until I determine that it is safe through community feedback to submit my own pull requests! I understand the temptations, but you will pose a risk to your community of Pixel Tensor users if you do not wait for community testing of this module. If you are faced with users who have complaints about audio quality issues, kindly direct them over to this GitHub repository where I will provide releases and updates for them to try out.

---

# Pixel Tensor Audio Decompressor

### Made and tested by JoshuaDoes on a Pixel 6 Pro.

---

## What does this module do?

This module will patch and systemlessly replace your mixer_paths.xml to make the following changes (if you choose the recommended preset):
- Raise the digital PCM volume from the stock value of 817 to 865.
- Lower the hardware amplifier gain from the stock value of 17 to 14.
- Set the boost peak current limit to the stock value of 3.50A. Some custom ROMs opt to lower it to reduce the heat output, at the expense of your phone's volume potential.

Using a newly integrated text menu engine, you may also choose from a variety of presets (PCM/AMP):
- 865/14: The recommended volume by my testing
- 865/13: The original volume from v1.0.0
- 865/12: Recommended by a Pixel 7 Pro (cheetah) user
- 817/17: The Google stock volume
- 913/15: Potentially unsafe! Triggers volume normalization when the volume spikes
- 913/20: Potentially even more unsafe than 913/15! The absolute maximum that can be set

More presets may be added in future updates, also depending on community feedback.

## What do these values even mean?

To start, let's consider Google's stock values of 817/17 (PCM/amp). The PCM volume has an estimated minimum of 800 and a known maximum of 913, whereas the amplifier gain has a known value range of 0 to 20. Think of it this way: They've decided to go with a low digital volume, and instead opted to amplify it to their desired loudness through the hardware amplifier.

If you aren't yet aware, PCM (pulse code modulation) is a digital representation of an analog signal, sampled at a fixed interval (for example, 44.1KHz or 48KHz). In order to lower the volume of a PCM signal, you must use math to compress its signal, thus losing some details from the source signal. Even if you use a hardware amplifier to gain more loudness out of this signal, you cannot possibly recover the lost details. This is precisely the problem with the audio quality on the Pixel 6 and Pixel 7 series of devices.

I've been trialing various increased values for months now using my Pixel 6 Pro, everything up to 913/15 and 913/20 included. But one small problem with these maximized values - as good as they can sound compared to my chosen values - is that Google has a safety feature built into their DSP (digital sound processor) which triggers volume normalization when the volume spikes. This is a speaker protection feature, because you should not expect your phone to be as capable as your laptop speakers might be. We do not want to disable this safety measure because we do not want our speakers to burn out from the heat production, or worse. The compromise afforded to us is that this speaker protection is sanely well configured. When given a max input volume of 913/20 and playing loud audio, it will only lower the speaker volume temporarily, for a short period of time. My resolve was to find a balanced set of values that maximize our processed digital PCM volume while introducing no artifacting, no clipping, and no distortion. There was some light clipping taking place with the amplifier gain being kept at 17, potentially due to the speaker protection feature, and so I opted to lower it too.

## How does this affect the audio quality?

The final result with selecting 865/14 is that we have minimized the compression of the PCM audio being sent to our speakers, allowing a wider range of detail. The noise floor has been lowered and the noise ceiling has been raised. You should be able to audibly comprehend better instrument separation, more depth to the notes being played in each instrument (including vocals, or voices and environmental sounds in videos), clarity across the audio spectrum, and a more pronounced bass. And due to the nature of raising the software volume, everything sounds just a little louder than before. This is before even considering something like tuning your EQ settings with an app like Wavelet, which should finally give you some more control over the sound of your audio due to the perceivable lack of detail compression.

## How did you determine these values?

I set the amplifier gain to a value of 0 when trying to determine the best PCM volume, starting at 900 because of the volume normalization at 913. I then tried to keep raising the amplifier gain by 1, doing my best to reach the loudness I perceived with the stock values. Every time I triggered the volume normalization, I would reset the amplifier gain back to 0 and then subtract 5 from the PCM volume, only to try raising the amplifier volume again. 865 turned out to be the first value I could use that wouldn't trigger it, but using any amplifier gain above 14 (including the stock 17) would introduce the aforementioned clipping issue. It didn't matter though, as 14 still allows the new PCM volume to exceed the loudness that we had before, without triggering safety features.

## What devices are supported by this module?

You must be running Android 13 or Android 14, regardless of whether or not you have installed a custom ROM. Future Android versions will require testing before allowing the installation to continue.

- Pixel 6      (oriole)
- Pixel 6 Pro  (raven)
- Pixel 6a     (bluejay)
- Pixel 7      (panther)
- Pixel 7 Pro  (cheetah)
- Pixel 7a     (lynx)
- Pixel Fold   (felix)
- Pixel Tablet (tangorpro)

## How do I test the differences?

That's the easy part! Go listen to your favorite music, watch your favorite YouTubers, play your favorite games, and listen to all the sounds of the crispy videos you've recorded. The only clear way to test this module is to listen to the things you're already used to listening to using your phone's speakers. If you want to do comparisons, use the "Choose a preset volume ..." menu to switch between your chosen preset and the Google stock preset (817/17). You can now live test the volume by temporarily applying it until the next reboot!

## Downloads

Release downloads: [GitHub releases for Pixel Tensor Audio Decompressor](https://github.com/JoshuaDoes/pixel-tensor-audio-decompressor/releases)

Source code: [GitHub repository for Pixel Tensor Audio Decompressor](https://github.com/JoshuaDoes/pixel-tensor-audio-decompressor)

**NOTE:** This module is no longer upgradeable through Magisk's module list! If you see an update available, you must visit the GitHub Releases tab on this repository in order to download and manually install it. It is HIGHLY recommended to install it ASAP. It may be a minor bug, or it may be something huge that I have overlooked. I want your phone to be as safe as possible when using your speakers.

## Support and Feedback

- [XDA](https://forum.xda-developers.com/t/mod-pixel-tensor-audio-decompressor.4610051/)

Old XDA threads that are currently closed by a moderator:

- [Pixel 6](https://forum.xda-developers.com/t/mod-pixel-tensor-audio-decompressor.4610053/)
- [Pixel 6 Pro](https://forum.xda-developers.com/t/mod-pixel-tensor-audio-decompressor.4610051/)
- [Pixel 6a](https://forum.xda-developers.com/t/mod-pixel-tensor-audio-decompressor.4610055/)
- [Pixel 7](https://forum.xda-developers.com/t/mod-pixel-tensor-audio-decompressor.4610069/)
- [Pixel 7 Pro](https://forum.xda-developers.com/t/mod-pixel-tensor-audio-decompressor.4610059/)
- [Pixel 7a](https://forum.xda-developers.com/t/mod-pixel-tensor-audio-decompressor.4610071/)
- [Pixel Fold](https://forum.xda-developers.com/t/mod-pixel-tensor-audio-decompressor.4610205/)
- [Pixel Tablet](https://forum.xda-developers.com/t/mod-pixel-tensor-audio-decompressor.4610201/)

Discord: @joshuadoes

Telegram: [JoshuaDoes](https://t.me/JoshuaDoes)

Twitter: [@TheNotesOfJosh](https://twitter.com/TheNotesOfJosh)

## Donations

If you enjoy the work I've done here, or any other projects I work on, please consider donating and leave a note with the projects that you used!

PayPal: [JoshuaDoes](https://paypal.me/JoshuaDoes)

Patreon: [JoshuaDoes](https://patreon.com/JoshuaDoes)

CashApp: $JoshuaDoes

Chime: $JoshuaDoes

Venmo: @JoshuaDoes
