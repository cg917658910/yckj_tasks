<?php
namespace app;

use think\db\exception\DataNotFoundException;
use think\db\exception\ModelNotFoundException;
use think\exception\Handle;
use think\exception\HttpException;
use think\exception\HttpResponseException;
use think\exception\ValidateException;
use think\Response;
use think\response\Json;
use Throwable;

/**
 * 应用异常处理类
 */
class ExceptionHandle extends Handle
{
    /**
     * 不需要记录信息（日志）的异常类列表
     * @var array
     */
    protected $ignoreReport = [
        HttpException::class,
        HttpResponseException::class,
        ModelNotFoundException::class,
        DataNotFoundException::class,
        ValidateException::class,
    ];

    /**
     * 记录异常信息（包括日志或者其它方式记录）
     *
     * @access public
     * @param  Throwable $exception
     * @return void
     */
    public function report(Throwable $exception): void
    {
        // 使用内置的方式记录异常日志
        parent::report($exception);
    }

    /**
     * Render an exception into an HTTP response.
     *
     * @access public
     * @param \think\Request   $request
     * @param Throwable $e
     * @return Response
     */
    public function render($request, Throwable $e): Response
    {
        if ($e instanceof HttpResponseException) {
            return $e->getResponse();
        }

        if ($this->shouldReturnJson($request)) {
            $message = $e->getMessage();
            $httpStatus = 500;
            $code = 1;

            if ($e instanceof ValidateException) {
                $httpStatus = 422;
            } elseif ($e instanceof HttpException) {
                $httpStatus = $e->getStatusCode();
            } elseif ($e instanceof ModelNotFoundException || $e instanceof DataNotFoundException) {
                $httpStatus = 404;
                $message = '资源不存在';
            }

            if ($message === '' || $message === 'error') {
                $message = '系统错误';
            }

            return json_error($message, $code, [], $httpStatus);
        }

        return parent::render($request, $e);
    }

    private function shouldReturnJson($request): bool
    {
        $path = $request->pathinfo();
        if (str_starts_with($path, 'admin/') || str_starts_with($path, 'user/')) {
            return true;
        }
        return $request->isJson() || $request->acceptJson();
    }
}
