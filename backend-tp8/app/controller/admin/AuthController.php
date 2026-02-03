<?php

declare(strict_types=1);

namespace app\controller\admin;

use app\controller\BaseController;
use app\service\AuthService;

class AuthController extends BaseController
{
    public function login()
    {
        $data = $this->request->only(['username', 'password'], 'post');
        try {
            validate(\app\validate\AdminLoginValidate::class)->check($data);
        } catch (\think\exception\ValidateException $e) {
            return json_error($e->getError());
        }

        $username = trim((string) $data['username']);
        $password = (string) $data['password'];

        try {
            $data = AuthService::adminLogin($username, $password);
        } catch (\Throwable $e) {
            return json_error($e->getMessage());
        }

        return json_success($data);
    }

    public function logout()
    {
        return json_success();
    }
}
