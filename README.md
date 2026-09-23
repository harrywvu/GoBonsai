# GoBonsai ⛩️🦫

GoBonsai is a Go script that procedurally generates ASCII art trees in your own terminal.
This project is greatly inspired by [PyBonsai](https://github.com/Ben-Edwards44/PyBonsai) which in itself is an ancestor to [cBonsai](https://github.com/mhzawadi/homebrew-cbonsai).

## Usage

```
gobonsai [OPTION]...

GoBonsai procedurally generates ASCII art trees in your terminal.

OPTIONS:
    -h, --help            display help
        --version         display version

    -s, --seed            seed for the random number generator

    -i, --instant         instant mode: display finished tree immediately
    -w, --wait            time delay between drawing characters when not in instant mode [default 0]

    -c, --branch-chars    string of chars randomly chosen for branches [default "~;:="]
    -C, --leaf-chars      string of chars randomly chosen for leaves [default "&%#@"]

    -x, --width           maximum width of the tree [default 80]
    -y, --height          maximum height of the tree [default 25]

    -t, --type            tree type: integer between 0 and 3 inclusive [default random]
    -S, --start-len       length of the root branch [default 15]
    -L, --leaf-len        length of each leaf [default 4]
    -l, --layers          number of branch layers: more => more branches [default 8]
    -a, --angle           mean angle of branches to their parent, in degrees; more => more arched trees [default 40]

    -f, --fixed-window    do not allow window height to increase when tree grows off screen
```

## License

[MIT License](https://opensource.org/license/mit)

TL;DR:
Do whatever you want with it. Just don't be a jerk.
