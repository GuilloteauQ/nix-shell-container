Output demo.gif

Set FontSize 20
Set Width 1920
Set Height 1080


Type "cat flake.nix"
Sleep 1s
Enter

Type "nix develop .#my-shell"
Sleep 1s
Enter

Type "pwd"
Sleep 1s
Enter

Type "exit"
Sleep 1s
Enter

Type "nano flake.nix"
Enter

Down 14
Tab 2
Type "containerize = true;"
Enter
Sleep 1s
Ctrl+X
Sleep 1s
Type "Y"
Sleep 1s
Enter

Type "nix develop .#my-shell"
Sleep 10s
Enter

Type "pwd"
Sleep 10s
Enter

Type "exit"
Sleep 10s
Enter
