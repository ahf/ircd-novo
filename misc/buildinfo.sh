#!/bin/sh
# vim: set sw=4 sts=4 et tw=80 :

# Copyright (c) 2011 Alexander Færøy <ahf@0x90.dk>
# All rights reserved.
#
# Redistribution and use in source and binary forms, with or without
# modification, are permitted provided that the following conditions are met:
#
# * Redistributions of source code must retain the above copyright notice, this
#   list of conditions and the following disclaimer.
#
# * Redistributions in binary form must reproduce the above copyright notice,
#   this list of conditions and the following disclaimer in the documentation
#   and/or other materials provided with the distribution.
#
# THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
# ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
# WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
# DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
# FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
# DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
# SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
# CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
# OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
# OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package=
misc_dir=$(dirname $0)

while [ $# -gt 0 ] ; do
    case $1 in
        --package)
            if [ $# -gt 1 ] ; then
                package=$2
                shift
            else
                echo "Please pass --package [package name]"
                exit 1
            fi
            ;;
        *)
            echo "Unknown parameter: $1"
            exit 1
            ;;
    esac
    shift
done

if [ -z $package ] ; then
    echo "Please pass --package [package name]"
    exit 1
fi

cat $misc_dir/generated-file.go

GITHEAD=`git describe 2>/dev/null`
if test -z "${GITHEAD}" ; then
    GITHEAD=`git rev-parse HEAD`
fi

if test -n "`git diff-index -m --name-only HEAD`" ; then
    GITHEAD="${GITHEAD}-dirty"
fi

echo
echo "package $package"
echo
echo "const ("
echo "    BUILDUSER = \"$(whoami)\""
echo "    BUILDHOST = \"$(hostname -f)\""
echo "    BUILDDATE = \"$(date +%Y-%m-%dT%H:%M:%S%z)\""
echo "    GITHEAD = \"$GITHEAD\""
echo ")"
