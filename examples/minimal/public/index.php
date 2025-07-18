<?php

use App\Sqlc\SQLite\QueriesImpl;

require_once __DIR__ . '/../vendor/autoload.php';

$pdo = new PDO('sqlite:' . __DIR__ . '/../data.db');
$sqlc = new QueriesImpl($pdo);

// $sqlc->createAuthor('some author');
$res = $sqlc->listAuthors();
foreach ($res as $author) {
    echo $author->authorId . ' ' . $author->name . PHP_EOL;
}
