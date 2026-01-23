# VHS Demo for brew-discover
# Usage: cd /path/to/brew-discover && vhs demo.tape

Output demo.gif

Set FontSize 14
Set Width 900
Set Height 550
Set Theme "Catppuccin Mocha"
Set TypingSpeed 50ms
Set Padding 20

Env PATH "/Users/akira/src/github.com/atani/brew-discover:/opt/homebrew/bin:/usr/bin:/bin"

Type "# brew-discover - Discover Homebrew packages"
Enter
Sleep 1s

Type "brew-discover top -n 5 --lang en"
Enter
Sleep 3s

Type "brew-discover random -n 3 --lang en"
Enter
Sleep 3s

Type "brew-discover browse --lang en"
Enter
Sleep 2s

Type "brew-discover info bat --lang en"
Enter
Sleep 3s

Type "brew-discover search editor -l 5 --lang en"
Enter
Sleep 3s
