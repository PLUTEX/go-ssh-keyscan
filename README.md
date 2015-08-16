# go-ssh-keyscan
Concurrent ssh-keyscan read hosts from STDIN and writes to STDOUT.

## Example

    % echo "www.maaret.de\nszuecs.net" | ./go-ssh-keyscan
    www.maaret.de ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBErFh4zfzhvotzqTN4tc+LMr/3ZENCOju1SLPQEthhugFws26NaqHbxnwtW+nCv+f1PAAs6RUuZe8lis4ggWxz8=
    szuecs.net ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBErFh4zfzhvotzqTN4tc+LMr/3ZENCOju1SLPQEthhugFws26NaqHbxnwtW+nCv+f1PAAs6RUuZe8lis4ggWxz8=
