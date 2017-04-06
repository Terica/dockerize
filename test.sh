#!/bin/bash
fail() {
	echo "$@"
	exit 1
}

grep -q docker /proc/$$/cgroup 2>/dev/null
if [ $? -ne 0 ]; then
	docker run -i --rm -v /var/run/docker.sock:/var/run/docker.sock -v /usr/bin/docker:/usr/bin/docker:ro -v ${PWD}:${PWD} -w ${PWD} ubuntu:14.04 "$0" "$@"
	exit
fi
uname -a
test -S /var/run/docker.sock || fail "You need to pass in docker and the docker socket."
version=$(docker version 2>&1) || fail -e "Docker is not functional in the container.\n$version"
cp dockerize /usr/local/bin
chmod +x /usr/local/bin/dockerize
test -e /usr/local/bin/golang && fail "golang shouldn't exist in this container to start with. no cheating."
dockerize install golang go || fail "dockerize golang go failed."
ls -li /usr/local/bin
test -e /usr/local/bin/golang || fail "we should have just created golang"
test -L /usr/local/bin/go || fail "we should have created a symlink go too"
[ $(ls -i /usr/local/bin/golang | cut -d' ' -f1) = $(ls -i /usr/local/bin/dockerize | cut -d' ' -f1) ] ||
	fail "Dockerize and golang should reference the same inode."
true
