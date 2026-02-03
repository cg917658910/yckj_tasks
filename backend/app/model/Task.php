<?php

declare(strict_types=1);

namespace app\model;

use think\Model;

class Task extends Model
{
    protected $table = 'tasks';
    protected $pk = 'id';
}
