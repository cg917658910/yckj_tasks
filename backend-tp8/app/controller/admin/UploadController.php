<?php

declare(strict_types=1);

namespace app\controller\admin;

use app\controller\BaseController;
use app\service\UploadService;

class UploadController extends BaseController
{
    public function image()
    {
        $file = $this->request->file('file');
        if (!$file) {
            return json_error('请选择文件');
        }

        try {
            $data = UploadService::saveImage($file);
        } catch (\Throwable $e) {
            return json_error($e->getMessage());
        }

        return json_success($data);
    }
}
