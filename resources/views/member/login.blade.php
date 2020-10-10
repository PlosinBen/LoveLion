@extends('basic')

@section('Nav', '')
@section('Footer', '')

@section('StyleSheet')
    <link rel="stylesheet" href="{{ asset('res/bootstrap-social-5.1.1.min.css') }}">
@endsection

@section('Content')
    <div class="row justify-content-center">
        <div class="col-10 col-md-6 col-lg-4 col-xl-3">
            <div class="card mt-5">
                <div class="card-header">
                    <h4 class="font-weight-light text-center my-1">Login</h4>
                </div>
                <div class="card-body">
                    <a href="{{ route('member.login.social', 'facebook') }}" class="btn btn-facebook btn-block">
                        <i class="fab fa-facebook-f fa-fw"></i>
                        Sign in with Facebook
                    </a>
                </div>
            </div>
        </div>
    </div>
@endsection
