<?php

declare(strict_types=1);

namespace app\service;

use think\File;
use think\facade\Config;
use \think\facade\Filesystem;

class UploadService
{
    public static function saveImage1(File $file): array
    {
        $uploadDir = root_path() . 'public/uploads';
        $maxSize = 50 * 1024 * 1024;
        $ext = strtolower((string) $file->getExtension());
        $allowed = ['jpg', 'jpeg', 'png', 'gif'];

        if ($file->getSize() > $maxSize) {
            throw new \RuntimeException('文件过大');
        }
        if (!in_array($ext, $allowed, true)) {
            //throw new \RuntimeException('不支持的文件类型');
        }

        $saveName = date('Ymd') . '/' . uniqid('img_', true) . '.' . $ext;
        $info = $file->move($uploadDir, $saveName);
        if (!$info) {
            throw new \RuntimeException('上传失败');
        }
        $relative = '/uploads/' . str_replace('\\', '/', $saveName);
        $url = request()->domain() . $relative;

        return [
            'url' => $url,
            'path' => $relative,
        ];
    }
    public static function saveImage(File $file){

        $file = request()->file('file');
        $hashRule = 'md5';
        $savename = request()->post('savename',$file->hashName($hashRule));
        $topic = 'cg';
        $disk = 'public';
        $diskConfig = Config::get('filesystem.disks.' . $disk);
        $savename = Filesystem::disk($disk)->putFileAs( $topic, $file, $savename);
  
        $defaultDiskUrl = '/storage/';
        $diskUrl = isset($diskConfig['url']) ? $diskConfig['url'] . '/' : $defaultDiskUrl;
        $url = request()->domain() .$diskUrl . $savename;
        $data = ["url" => $url,'path' => $savename];

        return ($data);
    
    }
}


