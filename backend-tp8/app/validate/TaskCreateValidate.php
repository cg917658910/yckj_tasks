<?php

declare(strict_types=1);

namespace app\validate;

use think\Validate;

class TaskCreateValidate extends Validate
{
    protected $rule = [
        'title' => 'require|max:100',
        'summary' => 'require|max:255',
        'detail' => 'require',
        'reward_points' => 'require|integer|gt:0',
    ];

    protected $message = [
        'title.require' => '请输入任务标题',
        'summary.require' => '请输入任务简介',
        'detail.require' => '请输入任务详情',
        'reward_points.require' => '请输入奖励积分',
    ];
}
