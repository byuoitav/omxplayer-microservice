# How to Install and Use the OMXPlayer-Microservice

1. Download and Flash the latest version of Raspberry Pi OS Lite to an SD Card
   * Note: You can chose to use the "with desktop and recommended software" version of Raspberry Pi OS instead, which will
   give you a full blown Desktop Environment as well as several other packages. This guide however
   will assume that the lite version has been used
2. Put the SD card into a Pi with an attached screen of some kind (such as the official 7" touchscreen) and an attached keyboard
   * Note: You may also chose to SSH into the pi instead of attaching a keyboard.
3. Log into the pi using the default image username (pi) and password (raspberry).
4. Update the OS to the latest packages
   * `sudo apt update && sudo apt upgrade`
5. Install OMXPlayer
   * `sudo apt install omxplayer`
6. Install i3, Chrome, and related dependencies
   * `sudo apt install xserver-xorg xinit i3 chromium-browser`
   * Note: You may chose a different window manager, desktop environment, or browser here.
   We use i3 by default as we will only be showing a webpage in chrome and i3 make that simple.
7. Create an i3 config file
   * Note: Here is an example config file, you may change whatever you like, though the very last line is what starts chrome
   and points it at the correct page.
   ```text
   # i3 config file (v4)
   #
   # Please see http://i3wm.org/docs/userguide.html for a complete reference!
   #
   # This config file uses keycodes (bindsym) and was written for the QWERTY
   # layout.
   #
   # To get a config file with the same key positions, but for your current
   # layout, use the i3-config-wizard
   #
   
   # Font for window titles. Will also be used by the bar unless a different font
   # is used in the bar {} block below.
   # This font is widely installed, provides lots of unicode glyphs, right-to-left
   # text rendering and scalability on retina/hidpi displays (thanks to pango).
   font pango:DejaVu Sans Mono 8
   # Before i3 v4.8, we used to recommend this one as the default:
   # font -misc-fixed-medium-r-normal--13-120-75-75-C-70-iso10646-1
   # The font above is very space-efficient, that is, it looks good, sharp and
   # clear in small sizes. However, its unicode glyph coverage is limited, the old
   # X core fonts rendering does not support right-to-left and this being a bitmap
   # font, it doesn’t scale on retina/hidpi displays.
   
   # use these keys for focus, movement, and resize directions when reaching for
   # the arrows is not convenient
   set $up l
   set $down k
   set $left j
   set $right semicolon
   
   # use Mouse+Mod4 to drag floating windows to their wanted position
   floating_modifier Mod4
   
   # start a terminal
   bindsym Mod4+Return exec i3-sensible-terminal
   
   # kill focused window
   bindsym Mod4+Shift+q kill
   
   # start dmenu (a program launcher)
   bindsym Mod4+d exec dmenu_run
   # There also is the (new) i3-dmenu-desktop which only displays applications
   # shipping a .desktop file. It is a wrapper around dmenu, so you need that
   # installed.
   # bindsym Mod4+d exec --no-startup-id i3-dmenu-desktop
   
   # change focus
   bindsym Mod4+$left focus left
   bindsym Mod4+$down focus down
   bindsym Mod4+$up focus up
   bindsym Mod4+$right focus right
   
   # alternatively, you can use the cursor keys:
   bindsym Mod4+Left focus left
   bindsym Mod4+Down focus down
   bindsym Mod4+Up focus up
   bindsym Mod4+Right focus right
   
   # move focused window
   bindsym Mod4+Shift+$left move left
   bindsym Mod4+Shift+$down move down
   bindsym Mod4+Shift+$up move up
   bindsym Mod4+Shift+$right move right
   
   # alternatively, you can use the cursor keys:
   bindsym Mod4+Shift+Left move left
   bindsym Mod4+Shift+Down move down
   bindsym Mod4+Shift+Up move up
   bindsym Mod4+Shift+Right move right
   
   # split in horizontal orientation
   bindsym Mod4+h split h
   
   # split in vertical orientation
   bindsym Mod4+v split v
   
   # enter fullscreen mode for the focused container
   bindsym Mod4+f fullscreen
   
   # change container layout (stacked, tabbed, toggle split)
   bindsym Mod4+s layout stacking
   bindsym Mod4+w layout tabbed
   bindsym Mod4+e layout toggle split
   
   # toggle tiling / floating
   bindsym Mod4+Shift+space floating toggle
   
   # change focus between tiling / floating windows
   bindsym Mod4+space focus mode_toggle
   
   # focus the parent container
   bindsym Mod4+a focus parent
   
   # focus the child container
   #bindsym Mod4+d focus child
   
   # move the currently focused window to the scratchpad
   bindsym Mod4+Shift+minus move scratchpad
   
   # Show the next scratchpad window or hide the focused scratchpad window.
   # If there are multiple scratchpad windows, this command cycles through them.
   bindsym Mod4+minus scratchpad show
   
   # switch to workspace
   bindsym Mod4+1 workspace 1
   bindsym Mod4+2 workspace 2
   bindsym Mod4+3 workspace 3
   bindsym Mod4+4 workspace 4
   bindsym Mod4+5 workspace 5
   bindsym Mod4+6 workspace 6
   bindsym Mod4+7 workspace 7
   bindsym Mod4+8 workspace 8
   bindsym Mod4+9 workspace 9
   bindsym Mod4+0 workspace 10
   
   # move focused container to workspace
   bindsym Mod4+Shift+1 move container to workspace 1
   bindsym Mod4+Shift+2 move container to workspace 2
   bindsym Mod4+Shift+3 move container to workspace 3
   bindsym Mod4+Shift+4 move container to workspace 4
   bindsym Mod4+Shift+5 move container to workspace 5
   bindsym Mod4+Shift+6 move container to workspace 6
   bindsym Mod4+Shift+7 move container to workspace 7
   bindsym Mod4+Shift+8 move container to workspace 8
   bindsym Mod4+Shift+9 move container to workspace 9
   bindsym Mod4+Shift+0 move container to workspace 10
   
   # reload the configuration file
   bindsym Mod4+Shift+c reload
   # restart i3 inplace (preserves your layout/session, can be used to upgrade i3)
   bindsym Mod4+Shift+r restart
   # exit i3 (logs you out of your X session)
   bindsym Mod4+Shift+e exec "i3-nagbar -t warning -m 'You pressed the exit shortcut. Do you really want to exit i3? This will end your X session.' -b 'Yes, exit i3' 'i3-msg exit'"
   
   # resize window (you can also use the mouse for that)
   mode "resize" {
           # These bindings trigger as soon as you enter the resize mode
   
           # Pressing left will shrink the window’s width.
           # Pressing right will grow the window’s width.
           # Pressing up will shrink the window’s height.
           # Pressing down will grow the window’s height.
           bindsym $left       resize shrink width 10 px or 10 ppt
           bindsym $down       resize grow height 10 px or 10 ppt
           bindsym $up         resize shrink height 10 px or 10 ppt
           bindsym $right      resize grow width 10 px or 10 ppt
   
           # same bindings, but for the arrow keys
           bindsym Left        resize shrink width 10 px or 10 ppt
           bindsym Down        resize grow height 10 px or 10 ppt
           bindsym Up          resize shrink height 10 px or 10 ppt
           bindsym Right       resize grow width 10 px or 10 ppt
   
           # back to normal: Enter or Escape
           bindsym Return mode "default"
           bindsym Escape mode "default"
   }
   
   bindsym Mod4+r mode "resize"
   
   # Start i3bar to display a workspace bar (plus the system information i3status
   # finds out, if available)
   bar {
           status_command i3status
   }
   
   # Open Chromium on boot in Kiosk mode and set it to the main URL for Control Panel
   exec chromium-browser --kiosk --incognito --disable-pinch --disable-session-crashed-bubble --overscroll-history-navigation=0 --disable-infobars --simulate-outdated-no-au='Tue, 31 Dec 2099 23:59:59 GMT' http://localhost:8032/
   ```

8. Setup Xinit to start i3
   * Write the following to ~/.xinitrc
   ```text
   #!/usr/bin/env bash

   screenshutoff
   exec i3
   ```

9. Enable autologin for the pi user
   * Note: You can choose another user here if you wish, though you will need to make sure to setup the user correctly
   * Write the following tile to `/etc/systemd/system/getty@tty1.service.d/autologin.conf`
   ```text
   [Service]
   ExecStart=
   ExecStart=-/sbin/agetty --autologin pi --noclear %I 38400 linux
   ```
   
10. Build the OMXPlayer-Microservice. This can be done in two ways:
	1. Build the service on your local device and then copy the build output to the pi (fastest method)
		1. Make sure you have the following dependencies installed:
			* Git
			* Golang
			* Make
			* Tar
		2. Clone the repository to a location on your local device
			* `git clone https://github.com/byuoitav/omxplayer-microservice.git`
		3. `cd` into the directory
		4. Run `make build`
		5. Copy the tarball located at `./dist/omxplayer-microservice.tar.gz` to the pi
			* The easiest way to do this is either through `scp` or using a jump drive
		6. Extract the tarball on the pi in the desired location
			* `tar -xzvf omxplayer-microservice.tar.gz`
	 2. Build the service on the pi itself
		1. Make sure you install the required dependencies:
			* `sudo apt install git golang make`
		2. Clone the repository to a location on the pi
			* `git clone https://github.com/byuoitav/omxplayer-microservice.git`
		3. `cd` into the directory
		4. Run `make build`
		5. Move the build output to the desired location
			* `cp dist/omxplayer-microservice.tar.gz {desired location}`
		6. Extract the tarball 
			* `tar -xzvf omxplayer-microservice.tar.gz`

11. Install a systemd service to start omxplayer-microservice on boot (and keep it running)
	* Write the following file to `/etc/systemd/system/omxplayer-microservice.service`
	```text
	[Unit]
	Description=OMX Player Microservice

	[Service]
	WorkingDirectory={Location of untar-ed directory}/files
	ExecStart={Location of untar-ed directory}/omxplayer-microservice
	Restart=on-failure
	Environment="CACHE_DATABASE_LOCATION={Location of untar-ed directory}/cache.db"
	Environment="CONTROL_CONFIG_PATH={Location of untar-ed directory}/control-config.json"
	Environment="OMXPLAYER_DISPLAY=2"

	[Install]
	WantedBy=default.target
	```
	* Note: Please make sure to replace the file paths with the path where you have placed the files

12. Define the possible streams
	* Write the control-config file to the location where you untar-ed the files (`{location}/control-config.json`)
	* Here is an example:
	```json
	{
		"streams": [
			{
				"name": "First Stream",
				"url": "https://example.com/stream"
			},
			{
				"name": "Second Stream",
				"url": "udp://server/stream"
			}
		]
	}
	```

13. Restart the pi
	* When the pi restarts, it should:
		* Start the OMXPlayer-Microservice
		* Start i3
		* Launch a chrome window pointed at the microservice's webpage

14. Pick a stream to start
	



