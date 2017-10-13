# go-utils
go-utils is a set of utility libraries created to handle common use cases
for the type of projects I tend to write. Where a package has a Config struct
that struct has json tags for scanning in configuration.

## fs
Provides types and functions for replicating Linux commands like mount and df
and getting file ownership information like owner uid, username.

## input
Provides a convenience function for reading a single line input from the CLI

## json
Provides two packages, file for scanning/writing json text files and web for
wrapping json data and writing it to a http.ResponseWriter

## ldap
Wraps ldap connection, entry requests and entry updates for very general usage

## regex
Provides a library of common regexes for CLI input validation, matching html tags etc