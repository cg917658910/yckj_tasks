<?php

declare(strict_types=1);

namespace app\validate;

use think\Validate;

class WithdrawalApplyValidate extends Validate
{
    protected $rule = [
        'amount' => 'require|float|gt:0',
    ];

    protected $message = [
        'amount.require' => '请输入提现金额',
    ];
}
