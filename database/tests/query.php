<?php

mysqli_report(MYSQLI_REPORT_ERROR | MYSQLI_REPORT_STRICT);
$mysqli = new mysqli('localhost', 'root', 'root', 'test');

// $stmt = $mysqli->prepare("INSERT INTO CountryLanguage VALUES (?, ?, ?, ?)");
// $stmt->bind_param('sssd', $code, $language, $official, $percent);
//
// $code = 'DEU';
// $language = 'Bavarian';
// $official = "F";
// $percent = 11.2;
//
// $stmt->execute();
//
// printf("%d row inserted.\n", $stmt->affected_rows);


// $mysqli->query("DELETE FROM CountryLanguage WHERE Language='Bavarian'");
// printf("%d row deleted.\n", $mysqli->affected_rows);


$stmt = $mysqli->prepare("SELECT `age` FROM `user` WHERE `id` = ?");

// $stmt->bind_param('ss', ...['DEU', 'POL']);
$stmt->bind_param('i', ...[2]);
$stmt->execute();
$stmt->store_result();  // 否则 num_rows() 结果为 0

printf("%d rows found.\n", $stmt->num_rows());

$meta = $stmt->result_metadata();
           
while ( $field = $meta->fetch_field() ) 
{
    $parameters[] = &$row[$field->name];
} 

call_user_func_array(array($stmt, 'bind_result'), $parameters);

$results = [];

while ( $stmt->fetch() ) 
{ 
    $x = array(); 
    foreach( $row as $key => $val ) 
    { 
        $x[$key] = $val; 
    } 
    $results[] = $x; 
}

$stmt->close();
$mysqli->close();

print_r($results);

var_dump($results[0]["age"]);    
