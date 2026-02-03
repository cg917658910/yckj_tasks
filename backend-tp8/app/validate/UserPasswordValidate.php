<?php

declare(strict_types=1);

namespace app\validate;

use think\Validate;

class UserPasswordValidate extends Validate
{
    protected $rule = [
        'old_password' => 'require|min:6|max:50',
        'new_password' => 'require|min:6|max:50',
    ];

    protected $message = [
        'old_password.require' => '请输入原密码',
        'new_password.require' => '请输入新密码',
    ];
}
