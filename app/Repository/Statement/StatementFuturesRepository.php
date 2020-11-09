<?php

namespace App\Repository\Statement;

use App\Models\Investment\StatementFutures;
use App\Repository\Repository;
use Illuminate\Database\Eloquent\Model;

class StatementFuturesRepository extends Repository
{
    public function __construct(StatementFutures $futuresStatement)
    {
        $this->Model = $futuresStatement;
    }

    public function insert($columns): Model
    {
        $columns['oversea_commitment'] = $columns['oversea_commitment'] ?? 0;
        $columns['note'] = $columns['note'] ?? '';

        return parent::insert($columns);
    }
}
