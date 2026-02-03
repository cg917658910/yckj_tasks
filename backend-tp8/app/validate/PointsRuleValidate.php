<?php

declare(strict_types=1);

namespace app\validate;

use think\Validate;

class PointsRuleValidate extends Validate
{
    protected $rule = [
        'exchange_rate' => 'require|integer|gt:0',
        'min_withdraw_amount' => 'require|float|gt:0',
        'register_bonus_points' => 'require|integer|egt:0',
    ];

    protected $message = [
        'exchange_rate.require' => '请填写兑换比例',
        'min_withdraw_amount.require' => '请填写最低提现金额',
        'register_bonus_points.require' => '请填写注册赠送积分',
    ];
}
