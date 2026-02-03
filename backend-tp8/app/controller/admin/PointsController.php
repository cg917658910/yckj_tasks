<?php

declare(strict_types=1);

namespace app\controller\admin;

use app\controller\BaseController;
use app\service\PointsService;
use think\facade\Db;

class PointsController extends BaseController
{
    public function rules()
    {
        $rule = PointsService::getRules();
        return json_success($rule->toArray());
    }

    public function updateRules()
    {
        $adminId = (int) ($this->request->admin['admin_id'] ?? 0);
        $data = [
            'exchange_rate' => (int) $this->request->post('exchange_rate', 10),
            'min_withdraw_amount' => (float) $this->request->post('min_withdraw_amount', 10),
            'register_bonus_points' => (int) $this->request->post('register_bonus_points', 10),
        ];
        try {
            validate(\app\validate\PointsRuleValidate::class)->check($data);
        } catch (\think\exception\ValidateException $e) {
            return json_error($e->getError());
        }

        try {
            $rule = PointsService::updateRules($adminId, $data);
        } catch (\Throwable $e) {
            return json_error($e->getMessage());
        }

        return json_success($rule->toArray());
    }

    public function logs()
    {
        $userId = (int) $this->request->get('user_id', 0);
        $type = trim((string) $this->request->get('type', ''));
        $keyword = trim((string) $this->request->get('keyword', ''));
        $page = (int) $this->request->get('page', 1);
        $pageSize = (int) $this->request->get('page_size', 10);

        $query = Db::name('points_logs')->alias('l')
            ->leftJoin('users u', 'l.user_id = u.id')
            ->field('l.*, u.username')
            ->order('l.id', 'desc');

        if ($userId > 0) {
            $query->where('l.user_id', $userId);
        }
        if ($type !== '') {
            $query->where('l.type', $type);
        }
        if ($keyword !== '') {
            $query->whereLike('u.username', '%' . $keyword . '%');
        }

        $pager = $query->paginate(['list_rows' => $pageSize, 'page' => $page])->toArray();

        return json_success([
            'list' => $pager['data'] ?? [],
            'total' => (int) ($pager['total'] ?? 0),
            'page' => (int) ($pager['current_page'] ?? $page),
            'page_size' => (int) ($pager['per_page'] ?? $pageSize),
        ]);
    }
}
