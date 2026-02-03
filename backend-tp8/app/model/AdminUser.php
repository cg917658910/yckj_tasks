<?php

declare(strict_types=1);

namespace app\model;

use think\Model;

class AdminUser extends Model
{
    protected $table = 'admin_users';
    protected $pk = 'id';
}
