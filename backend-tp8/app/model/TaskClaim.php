<?php

declare(strict_types=1);

namespace app\model;

use think\Model;

class TaskClaim extends Model
{
    protected $table = 'task_claims';
    protected $pk = 'id';
}
