<?php

use App\Sqlc\MySQL\QueriesImpl;
use Doctrine\DBAL\Connection;
use Symfony\Bundle\FrameworkBundle\Test\KernelTestCase;
use Symfony\Component\Uid\Uuid;

class PersistenceTest extends KernelTestCase
{
    public function test()
    {
        self::bootKernel();

        $container = static::getContainer();

        $connection = $container->get(Connection::class);

        $sub = new QueriesImpl($connection);

        $date = new \DateTimeImmutable();
        $sub->createBook(
            authorId: 3, isbn: "someISBN", bookType: "someType", uuidToBin: Uuid::v7()->toString(), yr: 2323, available: $date, tags: "testTag"
        );
        $books = $sub->bookByTags("testTag");

        self::assertEquals(1, count($books));

        self::assertEquals("someISBN", $books[0]->isbn);
    }
}