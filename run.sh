path=$1
gopkgs=$(find $path -mindepth 2 -maxdepth 2 -type d)
echo "" > output.log
for p in $gopkgs
do
#    echo $p
    ccname=`echo $p | tr '/' ' ' | awk '{print $8}'`
    resultpath="./"$ccname"_result.txt"
    echo start-check:$ccname
    echo start-check:$ccname >> output.log
#    ./cc-checker -path $p
    ./cc-checker -path $p >> output.log
#    echo "./" + $p + "_result"+.txt
done

m1=$(cat output.log | grep "time" | wc -l)
echo "time variable detected:" $m1 >> output.log

m1=$(cat output.log | grep "use global variable" | wc -l)
echo "global variable detected:" $m1 >> output.log

m1=$(cat output.log | grep "rand" | wc -l)
echo "rand variable detected:" $m1 >> output.log

m1=$(cat output.log | grep "range query map" | wc -l)
echo "range query map detected:" $m1 >> output.log

m1=$(cat output.log | grep "command line" | wc -l)
echo "command line variable detected:" $m1 >> output.log

m1=$(cat output.log | grep "external file" | wc -l)
echo "external file detected:" $m1 >> output.log

m1=$(cat output.log | grep "external web" | wc -l)
echo "external web variable detected:" $m1 >> output.log

m1=$(cat output.log | grep "external env" | wc -l)
echo "external env detected:" $m1 >> output.log

m1=$(cat output.log | grep "sink here" | wc -l)
echo "Non-determined risk detected:" $m1 >> output.log

m1=$(cat output.log | grep "cross channel" | wc -l)
echo "cross channel invoke detected:" $m1 >> output.log

m1=$(cat output.log | grep "read after write" | wc -l)
echo "read after write detected:" $m1 >> output.log

m1=$(cat output.log | grep "range query" | wc -l)
echo "range query phantom read detected:" $m1 >> output.log