<?php

declare(strict_types=1);

namespace app\validate;

use think\Validate;

class AdminLoginValidate extends Validate
{
    protected $rule = [
        'username' => 'require|min:2|max:50',
        'password' => 'require|min:6|max:50',
    ];

    protected $message = [
        'username.require' => '请输入用户名',
        'password.require' => '请输入密码',
    ];
}
