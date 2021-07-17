<?php
echo __dir__."\n";
exit;
$arr = array(0,1,2,3,4,5,6);
 
if(!($nShmID = shm_attach(ftok(__FILE__, 'i'), 1024))) {
    die("create shared memory failed.\n");
}

$nVarKey = 1;
if(!shm_put_var($nShmID, $nVarKey, $arr)) {
    die("failed to put var\n");
}
 
 
if(!($arr1 = shm_get_var($nShmID, $nVarKey))) {
    die("failed to get arr1\n");
}
var_dump($arr1);
array_pop($arr1);
 
 
if(!($arr2 = shm_get_var($nShmID, $nVarKey))) {
    die("failed to get arr2\n");
}
 
 
if ($arr != $arr2) {
    echo "get a copy\n";
} else {
    echo "get a reference\n";
}
var_dump($arr2);
 
 
if(!shm_remove($nShmID)) {
    die("failed to remove shared memory\n");
}



