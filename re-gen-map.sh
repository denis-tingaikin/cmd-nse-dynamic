
prefix=$(go env GOMODCACHE)
repo=$1

deps=$(grep -r 'NewServer' $prefix/$repo -l --include \*server\*.go --exclude \*nse_\* --exclude \*ns_\*  --exclude \*test\* | sort --unique)
for dep in $deps 
do
    dep=${dep#$prefix/$repo/pkg/}
    echo $dep
done