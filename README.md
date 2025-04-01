# ABOUT
Passfu is a command-line password manager. It is currently under development and doesn't actually work yet. I intend for this to be a portfolio piece that I will pin to my Github's profile page. 

# DESIGN ISSUES
I couldn't get SQLCipher working with Go, so my workaround was to build an Encrypt()
and Decrypt() function (via the `easycipher` package) and have users manually encrypt / decrypt
the Sqlite DB with that instead. 


The current workflow looks like this:
- User creates a new db with `newdb` (unencrypted)
- User does work with the rest of the program
- User manually encrypts db
- User manually decrypts db before they use it again.

Obviously this isn't ideal because someone might forget to encrypt / decrypt.

I'm looking at ways to make this easier, because having to do it manually like that is prone
to error.

- Maybe there's a way to get sqlite.Open() to work with a []byte containing the decrypted contents of the Sqlite file?
- If not sqlite.Open() maybe gorm.Open()?
- iouil.Tempfile() might be a better solution: store the decrypted contents in a tempfile that's deleted after the program finishes, that way the db stays encrypted without much effort on the user's part
