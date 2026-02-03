<?php

declare(strict_types=1);

namespace app\controller\user;

use app\controller\BaseController;
use app\service\AuthService;

class AuthController extends BaseController
{
    public function register()
    {
        $username = trim((string) $this->request->post('username', ''));
        $password = (string) $this->request->post('password', '');

        if ($username === '' || $password === '') {
            return json_error('参数不完整');
        }

        try {
            $data = AuthService::userRegister($username, $password);
        } catch (\Throwable $e) {
            return json_error($e->getMessage());
        }

        return json_success($data);
    }

    public function login()
    {
        $username = trim((string) $this->request->post('username', ''));
        $password = (string) $this->request->post('password', '');

        if ($username === '' || $password === '') {
            return json_error('参数不完整');
        }

        try {
            $data = AuthService::userLogin($username, $password);
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
