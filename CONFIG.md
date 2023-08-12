# Configuration
Type `st config` in the terminal to configure sio-tool.
You will then see a few options.


(Alternatively, you can also just edit the `~/.st/config` file.)


## Login
With this option, you can login to Codeforces, Szkopul, or Sio.
The tool will try to log you in to the specified website and give you an error if one occurs.


After a successful login, your username and password will be saved to a corresponding session file, so you don't have to login again.


## Add a template
You will have a step-by-step guide to creating a template (for the tool to know what language you are using, what commands are used to compile and run your program, and what is your default template to copy whenever you want to solve a new problem).


You can have multiple templates, sio-tool will automatically choose the one for your program based on the file extension.


### Select a language
Choose what programming language this template is supposed to use.
(For C++17 (which Szkopul and Sio use), choose the option `GNU G++17 7.3.0`).


### Template absolute path
This is the 'absolute' (so it starts either at the root '/' or in your home directory '~/')
path to your template (when you want to solve a new problem, sio-tool will copy the contents of this file and replace some placeholders (for example, your username or date and time)).
For example, if your template is in your `Desktop` folder and it is a file called `template.cpp`,
write `~/Desktop/template.cpp`


The file extension of the template will be added to the list of file extensions matched by the template (in our example, it is `.cpp` so, now every file ending with `.cpp` can be used with this template).


### Other suffixes
If you want your template to match different file extensions, you can specify that (for example, `cxx` or `cc` is more or less equivalent to `cpp`, so if you sometimes also use those, you can write `"cxx cc"`).


### Template alias
This is purely an arbitrarily chosen name; it might as well be a random string, but it is good to name it, for example, with the language your template uses (for our example, `cpp`).


### Scripts
You will be asked to specify 3 scripts (`before_script`, `script` and `after_script`) you can think of them as compilation, execution, and clean up.


You can use some placeholders:
```plain
$%path%$ Path to
source file (Excluding $%full%$, e.g., "/home/arapak/")
$%full%$ Full name of source file (e.g., "a.cpp")
$%file%$ Name of source file (Excluding suffix, e.g., "a")
$%rand%$ Random string with 8 characters (including "a-z" "0-9")
```


So example scripts for C++ can be:
```
before_script (compilation): g++ -std=c++17 -Wall $%full%$ -o $%file%$
script (execution): ./$%file%$
after_script (clean up): rm ./$%file%$
```
The before script compiles your program, for example `abc.cpp` to a binary file `a` (so after replacing placeholders, the command would look like this: `g++ -std=c++17 -Wall a.cpp -o a`) then the script will run your program `./a` it will supply it with an input and check the output, so you don't have to do anything more here, and the after_script will clean up the binary, so delete the `a` file (with the command `rm ./a`).


### Make it default
You will be asked if you want this to be your default template from now on. It is the one that is actually used when you start solving a new problem and need a new file, but all the other functionality still works. So if you have two templates, one for c++ and one for python, and the c++ one is the default, sio-tool will create c++ files, but when you create a python file on your own, you can still use commands like `st test` to test your solution on the example test cases.


## Delete a template
Here you can delete a malformed template or just one you don't want to use or want to replace.


## Set the default template
If you forgot to mark your template as default or just want to change your default template, you can do it here.


## Run `st gen` after `st parse`
So `st gen` is the command used to copy your default template and create a new file for a problem you want to solve, `st parse` is a command to parse a contest and get all sample test cases, it creates a folder for every task in the contest, and here you can say if you want to automatically run `st gen` in every one of these folders, so you don't have to do it manually for every problem.


## Set codeforces host domain
If, for example, Codeforces is temporarily down and you want to compete in a contest, there are some different hosts, like https://m1.codeforces.com/, you can change the current host here.


## Set proxy
If you want to use a proxy, you can specify it here.


## Set folders' names
For every website, sio-tool has specified a path where you solve problems for the given site, the default ones are `~/st/codeforces`, `~/st/sio-staszic`, `~/st/sio-mimuw` and `~/st/szkopul`.


Also, for every archive in Szkopul and every section in Codeforces, there are also folders (for example, for Codeforces contests and gym, the default folders are `~/st/codeforces/contest` and `~/st/codeforces/gym`)


If you want to change those, do it here.


## Set default naming
For stress-testing purposes, you can specify the naming scheme for the solution file, brute force solution file, and generator file.


## Set database path
Every problem you parse is saved to a local SQLite database. Here, you can specify where the database file should be located.


# Configure your shell


To have some additional features (running stress tests in a new terminal window and working goto command), you can also configure your shell by adding the script below to your `.bashrc` file (or `.zshrc file).


```bash
st() {
  # If you want to run the stress-test command in a separate terminal, uncomment next block
#   if [ "$1" = stress-test ]; then
# 	command="st $@"
# 	# Now uncomment the command, which is for your terminal of choice
#		# this is an example for gnome-termianal (the default terminal for gnome users)
# 	# gnome-terminal -q -- bash -c "$command; read line"
# 		# konsole (default KDE terminal)
# 	# konsole --hold -e "$command"
# 		# xfce4-terminal (default xfce terminal)
# 	# xfce4-terminal -H -e "$command"
# 		# xterm (another popular terminal)
# 	# xterm -hold -e "$command"
# 		# terminator
# 	# terminator -e "$command" -p hold
# 		# alacritty terminal
# 	# alacritty --hold -e st "$@"
#     return
#   fi

  if [ "$1" = db ] && [ "$2" = goto ]; then
	res=$(command st "$@")
	code=$?
	if [ "$code" = 0 ]; then
		cd $res
	else
		echo -e "$res"
	fi
	else
  	command st "$@"
  fi
}
```