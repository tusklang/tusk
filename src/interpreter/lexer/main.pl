#! usr/bin/perl

use strict;
use JSON;
use File::Slurper 'read_text';
use warnings;

my $file = "";

foreach my $line(<STDIN>) {
    chomp($line);
    $file = $file . "\n" . $line;
}

# remove first newline
$file =~ s/^.{2}//s;

my $keywords_unDecoded = read_text("./lexer/keywords.json");
my $keywords = decode_json($keywords_unDecoded)

print $keywords_unDec;
