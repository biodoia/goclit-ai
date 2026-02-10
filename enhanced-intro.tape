# VHS tape for enhanced intro animation demo
# Creates: enhanced-intro.gif

Output enhanced-intro.gif

Set FontSize 14
Set Width 900
Set Height 550
Set Theme "Dracula"
Set Padding 20
Set TypingSpeed 0ms

# Run the enhanced intro demo
Type "go run ./cmd/introdemo"
Enter

# Wait for the animation to complete (3.5 seconds + buffer)
Sleep 4s

# Done - animation exits automatically
