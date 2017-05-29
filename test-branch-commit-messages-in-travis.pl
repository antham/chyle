#!/bin/perl

`git ls-remote origin master` =~ /([a-f0-9]{40})/;

my $refHead = `git rev-parse HEAD`;
my $refTail = $1;

chomp($refHead);
chomp($refTail);

if ($refHead eq $refTail) {
    exit 0;
}

system "gommit check range $refTail $refHead";

if ($? > 0) {
    exit 1;
}
