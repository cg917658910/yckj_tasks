<?php

declare(strict_types=1);

namespace app\controller\user;

use app\controller\BaseController;
use app\service\WithdrawalService;
use think\facade\Db;

class WithdrawalController extends BaseController
{
    public function apply()
    {
        $userId = (int) ($this->request->user['user_id'] ?? 0);
        $data = $this->request->only(['amount'], 'post');
        try {
            validate(\app\validate\WithdrawalApplyValidate::class)->check($data);
        } catch (\think\exception\ValidateException $e) {
            return json_error($e->getError());
        }

        $amount = (float) $data['amount'];

        try {
            $withdrawal = WithdrawalService::apply($userId, $amount);
        } catch (\Throwable $e) {
            return json_error($e->getMessage());
        }

        return json_success($withdrawal->toArray());
    }

    public function index()
    {
        $userId = (int) ($this->request->user['user_id'] ?? 0);
        $list = Db::name('withdrawals')->where('user_id', $userId)->order('id', 'desc')->select()->toArray();
        return json_success(['list' => $list]);
    }
}
