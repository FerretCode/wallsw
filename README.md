# wallsw-go
A wallpaper switcher written in Golang for use with https://github.com/doisoundlikeababy/wallpaper

## usage
```
git clone https://github.com/ferretcode/wallsw ~/wallsw
export PATH=$PATH:/home/$USER/wallsw/bin
```

You can switch to a random wallpaper by running:
```
wallsw -d <dirname> --random
```

Or step through all your wallpapers:
```
wallsw -d <dirname>
```

## arguments
- `-d` The directory path to your wallpapers
- `--random` A flag to determine if you are fetching a random wallpaper or not
