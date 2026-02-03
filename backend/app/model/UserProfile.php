<?php

declare(strict_types=1);

namespace app\model;

use think\Model;

class UserProfile extends Model
{
    protected $table = 'user_profiles';
    protected $pk = 'id';
}
