<?php

$sqlStr = "Select * From user where id = 10";
$sqlStr = ltrim(substr($sqlStr, 0, 11), '(');
var_dump($sqlStr);    

$stmt = preg_split('/[\s]+/', $sqlStr, 2);
var_dump($stmt);    
$stmt = reset($stmt);
var_dump($stmt);    



