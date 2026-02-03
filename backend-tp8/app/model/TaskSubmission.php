<?php

declare(strict_types=1);

namespace app\model;

use think\Model;

class TaskSubmission extends Model
{
    protected $table = 'task_submissions';
    protected $pk = 'id';
}
