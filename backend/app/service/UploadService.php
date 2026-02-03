<?php

declare(strict_types=1);

namespace app\service;

use think\File;

class UploadService
{
    public static function saveImage(File $file): string
    {
        $uploadDir = root_path() . 'public/uploads';
        $info = $file->validate(['size' => 5 * 1024 * 1024, 'ext' => 'jpg,jpeg,png,gif'])->move($uploadDir);
        if (!$info) {
            throw new \RuntimeException('上传失败');
        }

        return '/uploads/' . str_replace('\\\', '/', $info->getSaveName());
    }
}
