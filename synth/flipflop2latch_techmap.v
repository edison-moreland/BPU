// map flip flops into dlatch based implementations

module \$_DLATCH_P_ (D, E, Q);
    wire [1023:0] _TECHMAP_DO_ = "simplemap; opt";
    input D, E;
    output Q;

    DLATCH DLP(.E(E), .D(D), .Q(Q));
endmodule

module \$_DFF_P_ (D, C, Q);
    wire [1023:0] _TECHMAP_DO_ = "simplemap; opt";
    input D, C;
    output Q;
    wire QQ;

    DLATCH DLP1(.E(~C), .D(D), .Q(QQ));
    DLATCH DLP2(.E(C), .D(QQ), .Q(Q));
endmodule

// module \$_DFF_PP0_ (D, C, R, Q);
//     wire [1023:0] _TECHMAP_DO_ = "simplemap; opt";
//     input D, C, R;
//     output reg Q;
//     wire QQ;

//     DLATCH DLP1(.E(~C), .D(D), .Q(QQ));
//     DLATCH DLP2(.E(C|R), .D(QQ&(~R)), .Q(Q));
// endmodule

module \$_ALDFF_PP_ (D, C, L, AD, Q);
    wire [1023:0] _TECHMAP_DO_ = "simplemap; opt";
    input D, C, L, AD;
    output reg Q;

    wire DD, CC, QQ;

    assign CC = L | C;
    assign DD = (L & AD) | ( ~L & C & D);

    DLATCH DLP1(.E(~CC), .D(DD), .Q(QQ));
    DLATCH DLP2(.E(CC), .D(QQ), .Q(Q));

endmodule

module \$_SR_PP_ (S, R, Q);
    wire [1023:0] _TECHMAP_DO_ = "simplemap; opt";
    input S, R;
    output Q;
    wire QQ;

    DLATCH DLP1(.E(~S), .D(1), .Q(QQ));
    DLATCH DLP2(.E(R | S), .D(QQ & ~R), .Q(Q));
endmodule