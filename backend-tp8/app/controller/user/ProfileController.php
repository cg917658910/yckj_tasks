<?php

declare(strict_types=1);

namespace app\controller\user;

use app\controller\BaseController;
use app\model\User;
use think\facade\Db;

class ProfileController extends BaseController
{
    public function info()
    {
        $userId = (int) ($this->request->user['user_id'] ?? 0);
        $data = Db::name('users')->alias('u')
            ->leftJoin('points_accounts p', 'u.id = p.user_id')
            ->leftJoin('user_profiles up', 'u.id = up.user_id')
            ->field('u.id,u.username,u.created_at,p.available_points,p.frozen_points,p.withdrawn_points,up.wechat_qr_url,up.total_withdrawn')
            ->where('u.id', $userId)
            ->find();

        return json_success($data ?: []);
    }

    public function changePassword()
    {
        $userId = (int) ($this->request->user['user_id'] ?? 0);
        $data = $this->request->only(['old_password', 'new_password'], 'post');
        try {
            validate(\app\validate\UserPasswordValidate::class)->check($data);
        } catch (\think\exception\ValidateException $e) {
            return json_error($e->getError());
        }

        $old = (string) $data['old_password'];
        $new = (string) $data['new_password'];

        $user = User::where('id', $userId)->find();
        if (!$user || !password_verify($old, $user->password_hash)) {
            return json_error('原密码错误');
        }

        $user->password_hash = password_hash($new, PASSWORD_BCRYPT);
        $user->save();

        return json_success();
    }

    public function updateWechatQr()
    {
        $userId = (int) ($this->request->user['user_id'] ?? 0);
        $url = trim((string) $this->request->post('wechat_qr_url', ''));
        if ($url === '') {
            return json_error('请上传二维码');
        }

        Db::name('user_profiles')->where('user_id', $userId)->update([
            'wechat_qr_url' => $url,
            'updated_at' => date('Y-m-d H:i:s'),
        ]);

        return json_success();
    }
}
