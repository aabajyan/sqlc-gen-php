<?php

namespace App\Controller;

use App\Sqlc\MySQL\QueriesImpl;
use Doctrine\DBAL\Connection;
use Doctrine\DBAL\Exception;
use Symfony\Bundle\FrameworkBundle\Controller\AbstractController;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\Routing\Attribute\Route;
use Symfony\Component\Uid\Uuid;

class TestController extends AbstractController
{
    public function __construct(private readonly Connection $connection)
    {
    }

    #[Route('/api/test', methods: ['GET'])]
    public function test(): JsonResponse
    {
        try {
            $id = (new QueriesImpl($this->connection))->createBook(
                authorId: 3, isbn: "someISBN", bookType: "someType", uuidToBin: Uuid::v7()->toString(), yr: 2323, available: new \DateTimeImmutable(), tags: "kausimausi"
            );
            var_dump($id);
        } catch (Exception $e) {
            var_dump($e);
        }

        return new JsonResponse();
    }
}