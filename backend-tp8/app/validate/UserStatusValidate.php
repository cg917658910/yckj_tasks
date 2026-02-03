<?php

declare(strict_types=1);

namespace app\validate;

use think\Validate;

class UserStatusValidate extends Validate
{
    protected $rule = [
        'status' => 'require|in:0,1',
    ];

    protected $message = [
        'status.require' => '请设置状态',
        'status.in' => '状态值不合法',
    ];
}
