<?php

declare(strict_types=1);

namespace app\controller;

use think\App;
use think\Request;

class BaseController
{
    protected Request $request;

    public function __construct(App $app)
    {
        $this->request = $app->request;
    }
}
