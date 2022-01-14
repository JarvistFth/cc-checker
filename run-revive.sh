path=$1
gopkgs=$(find $path -type f -name "*.go")
echo "" > revive_output.log
for p in $gopkgs
do
#    echo $p
    ccname=$p
#    resultpath="./"$ccname"_result.txt"
    echo start-check:$ccname
#    echo start-check:$ccname >> revive_output.log
    ./revive $p >> revive_output.log
#    echo "./" + $p + "_result"+.txt
done

m1=$(cat revive_output.log | grep "global variable detected" | wc -l)
echo "global variable detected:" $m1 >> revive_output.log

m2=$(cat revive_output.log | grep "should not use the following blacklisted import" | wc -l)
echo "blacklisted import:" $m2 >> revive_output.log

m2=$(cat revive_output.log | grep "rand" | wc -l)
echo "use rand:" $m2 >> revive_output.log

m2=$(cat revive_output.log | grep "blacklisted import: \"time\"" | wc -l)
echo "use time:" $m2 >> revive_output.log

m2=$(cat revive_output.log | grep "os" | wc -l)
echo "use os:" $m2 >> revive_output.log

m2=$(cat revive_output.log | grep "net" | wc -l)
echo "use net:" $m2 >> revive_output.log

m3=$(cat revive_output.log | grep "should not use range over map" | wc -l)
echo "range over map:" $m3 >> revive_output.log

m4=$(cat revive_output.log | grep "should not read after write" | wc -l)
echo "read after write:" $m4 >> revive_output.log

m5=$(cat revive_output.log | grep "data obtained from phantom reads" | wc -l)
echo "phantom reads:" $m5 >> revive_output.log

m6=$(cat revive_output.log | grep "should not use goroutines" | wc -l)
echo "goroutines:" $m6 >> revive_output.log




