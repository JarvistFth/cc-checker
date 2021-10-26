path=$1
gopkgs=$(find $path -mindepth 2 -maxdepth 2 -type d)
for p in $gopkgs
do
    echo $p
    echo "\n"
done