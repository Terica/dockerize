
# dockerize

Run commands from docker containers

## Usage

download dockerize and put it in a directory in your path.

run `dockerize upgrade` to create utility links and keep application up to date.

## Why?

On a desert island, you get to keep only a handful of items.

One of them should be dockerize (assuming another is docker itself)

```bash
newlappy:dockerize faye$ docker run -it --rm -v /tmp/dockerize:/bin/execwdve:ro -v /var/run/docker.sock:/tmp/docker.sock -e DOCKER_HOST=unix:///tmp/docker.sock -v /usr/bin/docker:/usr/bin/docker:ro -v $HOME:$HOME -w $PWD golang:latest /bin/bash

root@e71856affc9a:/Users/faye/src/docker/dockerize/go/dockerize# ln -s /bin/execwdve /bin/ruby
root@e71856affc9a:/Users/faye/src/docker/dockerize/go/dockerize# ln -s /bin/execwdve /bin/irb 
root@e71856affc9a:/Users/faye/src/docker/dockerize/go/dockerize# ln -s /bin/execwdve /bin/bundle
root@e71856affc9a:~# cd /Users/faye/src/adventure/
root@e71856affc9a:/Users/faye/src/adventure# ls
CODE_OF_CONDUCT.md  Gemfile.lock  README.md  adventure.gemspec	lib
Gemfile		    LICENSE.txt   Rakefile   bin		spec
root@e71856affc9a:/Users/faye/src/adventure# ruby -v
ruby 2.3.1p112 (2016-04-26 revision 54768) [x86_64-linux]
root@e71856affc9a:/Users/faye/src/adventure# bundle
Fetching gem metadata from https://rubygems.org/..........
Fetching version metadata from https://rubygems.org/.
Fetching rake 10.5.0
Installing rake 10.5.0
Using adventure 0.1.0 from source at `.`
Fetching rspec-support 3.5.0
Installing rspec-support 3.5.0
Fetching diff-lcs 1.3
Installing diff-lcs 1.3
Using bundler 1.15.1
Fetching rspec-core 3.5.4
Installing rspec-core 3.5.4
Fetching rspec-expectations 3.5.0
Installing rspec-expectations 3.5.0
Fetching rspec-mocks 3.5.0
Installing rspec-mocks 3.5.0
Fetching rspec 3.5.0
Installing rspec 3.5.0
Bundle complete! 4 Gemfile dependencies, 9 gems now installed.
Bundled gems are installed into /usr/local/bundle.
root@e71856affc9a:/Users/faye/src/adventure# bin/console 
Adding player to object list
Adding apartment_door to object list
Adding apartment_key to object list
Adding office_desk_top_drawer to object list
Adding office_desk to object list
Adding office_cubicle_entrance to object list
Adding office_cubicle to object list
irb(main):001:0> look
An office cubicle
A nondescript space just large enough for a desk, chair and
human inhabitant, surrounded on three and a half sides by
partitions too high to look over.
Contents:
  An office desk
=> ["An office desk"]
irb(main):002:0> quit
```

After all, ruby is what golang containers were made for!

And if you can make _that_ work then imagine what it's like setting up your new laptop.