<?php

declare(strict_types=1);

namespace app\model;

use think\Model;

class Withdrawal extends Model
{
    protected $table = 'withdrawals';
    protected $pk = 'id';
}
