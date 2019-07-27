#!/bin/bash
target_systems=( "linux" "windows" )
target_architectures=( "amd64" "amd64" )

if [ -d ./bin ]
then
  rm -rf ./bin
fi

for i in $(seq 0 $((${#target_systems[@]}-1)) )
do
  target_os=${target_systems[i]}
  target_platforms=(${target_architectures[i]})
  for j in $(seq 0 $((${#target_platforms[@]}-1)) )
  do
    target_plat=${target_platforms[j]}
    output_file="./bin/rxsm-$target_os-$target_plat"
    if [ $target_os == "windows" ]
    then
      output_file+=".exe"
    fi
    echo $target_os $target_plat
    env GOOS=$target_os GOARCH=$target_plat go build -o $output_file src/main/playq.go
    if [ $? -ne 0 ]
    then
      exit $?
    fi
  done
done
# go build src/main/playq.go
echo "All targets built successfully"
exit 0
