<?php

declare(strict_types=1);

namespace app\validate;

use think\Validate;

class TaskSubmitValidate extends Validate
{
    protected $rule = [
        'images' => 'require|array|min:1',
        'images.*' => 'require|max:255',
        'remark' => 'max:1000',
    ];

    protected $message = [
        'images.require' => '请上传截图',
        'images.array' => '截图格式错误',
    ];
}
